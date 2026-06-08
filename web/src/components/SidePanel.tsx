import ConversationList from './ConversationList';
import { Conversation } from '../api/client';

interface SidePanelProps {
  isOpen: boolean;
  conversations: Conversation[];
  currentId: string | null;
  onSelectConversation: (id: string) => void;
  onNewConversation: () => void;
  onDeleteConversation: (id: string) => void;
  onClose: () => void;
}

export default function SidePanel({
  isOpen, conversations, currentId,
  onSelectConversation, onNewConversation, onDeleteConversation, onClose,
}: SidePanelProps) {
  return (
    <>
      {isOpen && (
        <div className="fixed inset-0 bg-black/40 z-20" onClick={onClose} />
      )}
      <div className={`fixed left-0 top-0 bottom-0 w-64 bg-surface-card border-r border-border z-30
        transform transition-transform duration-200 ${isOpen ? 'translate-x-0' : '-translate-x-full'}`}>
        <div className="flex flex-col h-full p-4">
          <div className="flex items-center justify-between mb-4">
            <span className="font-semibold text-sm">🧠 U-Hermes</span>
            <span className="text-2xs text-text-muted">v0.1.0</span>
          </div>

          <button
            onClick={onNewConversation}
            className="w-full py-2.5 mb-4 rounded-lg text-sm font-medium transition-all"
            style={{ background: 'linear-gradient(135deg, #6366f1, #7c3aed)' }}
          >
            + 新对话
          </button>

          <div className="text-2xs uppercase text-text-muted mb-2 tracking-wider">最近对话</div>

          <ConversationList
            conversations={conversations}
            currentId={currentId}
            onSelect={(id) => { onSelectConversation(id); onClose(); }}
            onDelete={onDeleteConversation}
          />

          <div className="mt-auto pt-4 border-t border-border">
            <a href="/settings" className="block py-2 text-sm text-text-secondary hover:text-text-primary transition-colors">
              ⚙️ 设置
            </a>
          </div>

          <div className="mt-4 p-3 rounded-xl text-center border border-dashed border-border"
            style={{ background: 'linear-gradient(135deg, rgba(99,102,241,0.06), rgba(139,92,246,0.04))' }}>
            <p className="text-2xs text-text-secondary leading-relaxed">
              🔮 U-Hermes 会越用越聪明
            </p>
          </div>
        </div>
      </div>
    </>
  );
}
