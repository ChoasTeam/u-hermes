interface Props {
  icon: string;
  title: string;
  description: string;
  actions?: { label: string; onClick: () => void; primary?: boolean }[];
}

export default function ErrorCard({ icon, title, description, actions }: Props) {
  return (
    <div className="flex gap-3 p-4 rounded-xl border border-warning/20 bg-surface-card"
      style={{ borderLeft: '3px solid var(--warning, #f59e0b)' }}>
      <span className="text-2xl flex-shrink-0">{icon}</span>
      <div>
        <h4 className="text-sm font-semibold text-text-primary mb-1">{title}</h4>
        <p className="text-xs text-text-secondary leading-relaxed mb-3">{description}</p>
        {actions && (
          <div className="flex gap-2">
            {actions.map((action) => (
              <button
                key={action.label}
                onClick={action.onClick}
                className={`px-3 py-1.5 rounded-lg text-xs font-medium transition-all ${
                  action.primary
                    ? 'text-white'
                    : 'bg-surface-hover text-text-secondary hover:text-text-primary'
                }`}
                style={action.primary ? { background: 'linear-gradient(135deg, #6366f1, #7c3aed)' } : undefined}
              >
                {action.label}
              </button>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
