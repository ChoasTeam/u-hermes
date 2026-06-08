interface TopBarProps {
  modelName: string;
  modelConnected: boolean;
  onToggleSidebar: () => void;
}

export default function TopBar({ modelName, modelConnected, onToggleSidebar }: TopBarProps) {
  return (
    <header className="flex items-center justify-between px-5 py-3 border-b border-border"
      style={{ background: 'linear-gradient(180deg, #111113 0%, #0a0a0b 100%)' }}>
      <div className="flex items-center gap-3">
        <button onClick={onToggleSidebar} className="text-text-muted hover:text-text-primary transition-colors">
          <svg width="20" height="20" viewBox="0 0 20 20" fill="currentColor">
            <path d="M3 5h14M3 10h14M3 15h14" stroke="currentColor" strokeWidth="1.5" fill="none"/>
          </svg>
        </button>
        <div className="flex items-center gap-2">
          <div className="w-7 h-7 rounded-lg flex items-center justify-center text-sm"
            style={{ background: 'linear-gradient(135deg, #6366f1, #8b5cf6)' }}>
            🧠
          </div>
          <span className="font-semibold text-sm tracking-tight">U-Hermes</span>
          <span className="w-1.5 h-1.5 rounded-full"
            style={{
              background: modelConnected ? '#22c55e' : '#f59e0b',
              boxShadow: modelConnected ? '0 0 6px rgba(34,197,94,0.4)' : '0 0 6px rgba(245,158,11,0.4)',
            }}
          />
        </div>
      </div>
      <div className="flex items-center gap-2">
        <span className="text-xs px-2.5 py-1 rounded-md bg-surface-hover text-text-secondary">{modelName || '未配置'}</span>
        <a href="/settings" className="text-text-muted hover:text-text-primary transition-colors p-1">
          <svg width="18" height="18" viewBox="0 0 20 20" fill="currentColor">
            <path fillRule="evenodd" d="M11.49 3.17c-.38-1.56-2.6-1.56-2.98 0a1.532 1.532 0 01-2.286.948c-1.372-.836-2.942.734-2.106 2.106.54.886.061 2.042-.947 2.287-1.561.379-1.561 2.6 0 2.978a1.532 1.532 0 01.947 2.287c-.836 1.372.734 2.942 2.106 2.106a1.532 1.532 0 012.287.947c.379 1.561 2.6 1.561 2.978 0a1.533 1.533 0 012.287-.947c1.372.836 2.942-.734 2.106-2.106a1.533 1.533 0 01.947-2.287c1.561-.379 1.561-2.6 0-2.978a1.532 1.532 0 01-.947-2.287c.836-1.372-.734-2.942-2.106-2.106a1.532 1.532 0 01-2.287-.947zM10 13a3 3 0 100-6 3 3 0 000 6z" clipRule="evenodd"/>
          </svg>
        </a>
      </div>
    </header>
  );
}
