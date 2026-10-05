import { test, expect } from './fixtures';

test.describe('Backlog sorting', () => {
  test('beans are sorted by priority, then type, then title within each section', async ({
    beans,
    backlogPage
  }) => {
    beans.create('Todo Normal Task', { status: 'todo', priority: 'normal', type: 'task' });
    beans.create('Todo Normal Bug', { status: 'todo', priority: 'normal', type: 'bug' });
    beans.create('Todo High Feature', { status: 'todo', priority: 'high', type: 'feature' });
    beans.create('Draft Idea', { status: 'draft', priority: 'low', type: 'task' });
    beans.create('Stale Idea', { status: 'stale', priority: 'high', type: 'task' });

    await backlogPage.goto(5);

    const titles = await backlogPage.getBeanTitles();
    // Todo section first (sorted by priority, then type, then title), then Draft, then Stale
    expect(titles).toEqual([
      'Todo High Feature',
      'Todo Normal Bug',
      'Todo Normal Task',
      'Draft Idea',
      'Stale Idea'
    ]);
  });

  test('child beans are sorted by status, with stale between draft and completed', async ({
    beans,
    backlogPage
  }) => {
    const parentId = beans.create('Parent Epic', { status: 'todo', type: 'epic' });
    for (const status of ['completed', 'stale', 'draft', 'todo']) {
      const id = beans.create(`Child ${status}`, { status, type: 'task' });
      beans.run(['update', id, '--parent', parentId]);
    }

    await backlogPage.goto(5);

    await expect(async () => {
      const childTitles = await backlogPage
        .beanByTitle('Parent Epic')
        .locator('.bean-item [role="button"] > div > span.text-sm')
        .allTextContents();
      expect(childTitles.map((t) => t.trim())).toEqual([
        'Child todo',
        'Child draft',
        'Child stale',
        'Child completed'
      ]);
    }).toPass({ timeout: 5_000 });
  });

  test('list re-sorts when a bean priority changes on disk', async ({ beans, backlogPage }) => {
    const id1 = beans.create('Low Priority Task', {
      status: 'todo',
      priority: 'low',
      type: 'task'
    });
    beans.create('Normal Priority Task', { status: 'todo', priority: 'normal', type: 'task' });

    await backlogPage.goto(2);

    let titles = await backlogPage.getBeanTitles();
    expect(titles).toEqual(['Normal Priority Task', 'Low Priority Task']);

    // Change the low-priority bean to critical via CLI (filesystem change)
    beans.update(id1, { priority: 'critical' });

    // The bean should now appear first
    await expect(async () => {
      titles = await backlogPage.getBeanTitles();
      expect(titles).toEqual(['Low Priority Task', 'Normal Priority Task']);
    }).toPass({ timeout: 5_000 });
  });

  test('list re-sorts when a bean status changes on disk', async ({ beans, backlogPage }) => {
    const id1 = beans.create('A Draft Bean', { status: 'draft', type: 'task' });
    beans.create('B Todo Bean', { status: 'todo', type: 'task' });

    await backlogPage.goto(2);

    // Todo section comes before Draft section
    let titles = await backlogPage.getBeanTitles();
    expect(titles).toEqual(['B Todo Bean', 'A Draft Bean']);

    // Move the draft bean to todo
    beans.update(id1, { status: 'todo' });

    // Both are now todo, should sort by title
    await expect(async () => {
      titles = await backlogPage.getBeanTitles();
      expect(titles).toEqual(['A Draft Bean', 'B Todo Bean']);
    }).toPass({ timeout: 5_000 });
  });

  test('new bean appears in correct sorted position', async ({ beans, backlogPage }) => {
    beans.create('Zebra Task', { status: 'todo', type: 'task' });
    beans.create('Alpha Task', { status: 'todo', type: 'task' });

    await backlogPage.goto(2);

    let titles = await backlogPage.getBeanTitles();
    expect(titles).toEqual(['Alpha Task', 'Zebra Task']);

    // Create a new bean that should sort between the two
    beans.create('Middle Task', { status: 'todo', type: 'task' });

    await expect(async () => {
      titles = await backlogPage.getBeanTitles();
      expect(titles).toEqual(['Alpha Task', 'Middle Task', 'Zebra Task']);
    }).toPass({ timeout: 5_000 });
  });

  test('dragging a bean reorders it within the backlog', async ({ beans, backlogPage }) => {
    // Create beans with same status/priority/type so they sort by title
    beans.create('Alpha', { status: 'todo', type: 'task' });
    beans.create('Bravo', { status: 'todo', type: 'task' });
    beans.create('Charlie', { status: 'todo', type: 'task' });

    await backlogPage.goto(3);

    // Initial order: Alpha, Bravo, Charlie
    let titles = await backlogPage.getBeanTitles();
    expect(titles).toEqual(['Alpha', 'Bravo', 'Charlie']);

    // Drag Charlie above Alpha
    await backlogPage.dragBean('Charlie', 'Alpha', 'above');

    // New order: Charlie, Alpha, Bravo
    await expect(async () => {
      titles = await backlogPage.getBeanTitles();
      expect(titles).toEqual(['Charlie', 'Alpha', 'Bravo']);
    }).toPass({ timeout: 5_000 });
  });

  test('dragging a bean onto another reparents it and persists', async ({
    beans,
    backlogPage,
    page
  }) => {
    beans.create('Parent Bean', { status: 'todo', type: 'feature' });
    beans.create('Child Bean', { status: 'todo', type: 'task' });

    await backlogPage.goto(2);

    // Both are top-level initially
    let titles = await backlogPage.getBeanTitles();
    expect(titles).toContain('Parent Bean');
    expect(titles).toContain('Child Bean');

    // Drag Child Bean onto Parent Bean to reparent it
    await backlogPage.dragBean('Child Bean', 'Parent Bean', 'onto');

    // Child Bean should now be nested under Parent Bean
    await expect(async () => {
      const parentItem = backlogPage.beanByTitle('Parent Bean');
      const nestedChild = parentItem.locator('.bean-item', { hasText: 'Child Bean' });
      await expect(nestedChild).toBeVisible();
    }).toPass({ timeout: 5_000 });

    // Reload the page to verify persistence
    await page.reload();
    await backlogPage.waitForBean('Parent Bean');

    // After reload, Child Bean should still be nested under Parent Bean
    await expect(async () => {
      const parentItem = backlogPage.beanByTitle('Parent Bean');
      const nestedChild = parentItem.locator('.bean-item', { hasText: 'Child Bean' });
      await expect(nestedChild).toBeVisible();
    }).toPass({ timeout: 5_000 });
  });

  test('dragging a bean into a specific position within another parent works', async ({
    beans,
    backlogPage,
    page
  }) => {
    // Create a feature with 3 child tasks
    const parentId = beans.create('Parent Feature', { status: 'todo', type: 'feature' });
    const childA = beans.create('Alpha Child', { status: 'todo', type: 'task' });
    const childB = beans.create('Bravo Child', { status: 'todo', type: 'task' });
    const childC = beans.create('Charlie Child', { status: 'todo', type: 'task' });

    // Set parent relationships via CLI
    beans.run(['update', childA, '--parent', parentId]);
    beans.run(['update', childB, '--parent', parentId]);
    beans.run(['update', childC, '--parent', parentId]);

    // Create a top-level task to drag in
    beans.create('Interloper', { status: 'todo', type: 'task' });

    await backlogPage.goto(5); // Parent Feature + 3 children + Interloper

    // Drag Interloper above Bravo Child (between Alpha and Bravo)
    await backlogPage.dragBean('Interloper', 'Bravo Child', 'above');

    // Verify: Interloper should be between Alpha Child and Bravo Child
    await expect(async () => {
      const parentItem = backlogPage.beanByTitle('Parent Feature');
      const childTitles = await parentItem
        .locator('.bean-item [role="button"] > div > span.text-sm')
        .allTextContents();
      const trimmed = childTitles.map((t) => t.trim());
      expect(trimmed).toEqual(['Alpha Child', 'Interloper', 'Bravo Child', 'Charlie Child']);
    }).toPass({ timeout: 5_000 });

    // Verify persistence
    await page.reload();
    await backlogPage.waitForBean('Parent Feature');

    await expect(async () => {
      const parentItem = backlogPage.beanByTitle('Parent Feature');
      const childTitles = await parentItem
        .locator('.bean-item [role="button"] > div > span.text-sm')
        .allTextContents();
      const trimmed = childTitles.map((t) => t.trim());
      expect(trimmed).toEqual(['Alpha Child', 'Interloper', 'Bravo Child', 'Charlie Child']);
    }).toPass({ timeout: 5_000 });
  });

  test('deleted bean disappears from list', async ({ beans, backlogPage }) => {
    const id1 = beans.create('Bean To Delete', { status: 'todo', type: 'task' });
    beans.create('Bean To Keep', { status: 'todo', type: 'task' });

    await backlogPage.goto(2);

    // Delete the bean via CLI
    beans.run(['delete', '--force', id1]);

    await backlogPage.waitForBeanGone('Bean To Delete');
    expect(await backlogPage.count()).toBe(1);
  });
});
