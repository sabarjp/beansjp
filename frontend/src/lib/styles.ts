/** Shared color/style maps for bean badges, borders, and indicators. */

export const statusColors: Record<string, string> = {
  draft: 'bg-status-draft-bg text-status-draft-text',
  todo: 'bg-status-todo-bg text-status-todo-text',
  'in-progress': 'bg-status-in-progress-bg text-status-in-progress-text',
  stale: 'bg-status-stale-bg text-status-stale-text',
  completed: 'bg-status-completed-bg text-status-completed-text',
  scrapped: 'bg-status-scrapped-bg text-status-scrapped-text'
};

export const typeColors: Record<string, string> = {
  milestone: 'bg-type-milestone-bg text-type-milestone-text',
  epic: 'bg-type-epic-bg text-type-epic-text',
  feature: 'bg-type-feature-bg text-type-feature-text',
  bug: 'bg-type-bug-bg text-type-bug-text',
  task: 'bg-type-task-bg text-type-task-text'
};

export const typeBorders: Record<string, string> = {
  milestone: 'border-l-type-milestone-border',
  epic: 'border-l-type-epic-border',
  feature: 'border-l-type-feature-border',
  bug: 'border-l-type-bug-border',
  task: 'border-l-type-task-border'
};

export const priorityColors: Record<string, string> = {
  critical: 'border-danger text-danger',
  high: 'border-warning text-warning',
  normal: 'border-border text-text-muted',
  low: 'border-border text-text-muted opacity-60',
  deferred: 'border-border text-text-muted opacity-40'
};

export const priorityIndicators: Record<string, string> = {
  critical: 'text-danger',
  high: 'text-warning',
  low: 'text-text-faint',
  deferred: 'text-text-faint opacity-60'
};
