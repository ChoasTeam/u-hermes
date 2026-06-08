import { Conversation } from '../api/client';

interface Props {
  conversations: Conversation[];
  currentId: string | null;
  onSelect: (id: string) => void;
  onDelete: (id: string) => void;
}

export default function ConversationList({ conversations, currentId, onSelect, onDelete }: Props) {
  if (conversations.length === 0) {
    return <p className="text-xs text-text-muted text-center py-8">暂无对话记录</p>;
  }

  return (
    <div className="flex-1 overflow-y-auto space-y-0.5">
      {conversations.map((conv) => (
        <div
          key={conv.id}
          onClick={() => onSelect(conv.id)}
          className={`group flex items-center justify-between px-2.5 py-2 rounded-lg cursor-pointer text-sm transition-colors
            ${conv.id === currentId ? 'bg-surface-hover text-text-primary' : 'text-text-secondary hover:bg-surface-hover hover:text-text-primary'}`}
        >
          <span className="truncate flex-1">📝 {conv.title || '新对话'}</span>
          <button
            onClick={(e) => { e.stopPropagation(); onDelete(conv.id); }}
            className="opacity-0 group-hover:opacity-100 text-text-muted hover:text-danger transition-all text-xs px-1"
          >
            🗑️
          </button>
        </div>
      ))}
    </div>
  );
}
