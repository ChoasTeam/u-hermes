import ReactMarkdown from 'react-markdown';
import { Message } from '../api/client';

interface Props {
  message: Message;
  isStreaming?: boolean;
}

export default function MessageBubble({ message, isStreaming }: Props) {
  const isUser = message.role === 'user';

  return (
    <div className={`flex ${isUser ? 'justify-end' : 'justify-start'} mb-4`}>
      <div className={`max-w-[72%] ${isUser ? 'order-1' : ''}`}>
        <div className={`flex items-center gap-1.5 mb-1 ${isUser ? 'justify-end' : ''}`}>
          {isUser ? (
            <>
              <span className="text-2xs text-text-muted uppercase tracking-wider">你</span>
              <div className="w-4 h-4 rounded-full flex items-center justify-center text-2xs text-white font-bold"
                style={{ background: 'linear-gradient(135deg, #6366f1, #8b5cf6)' }}>
                你
              </div>
            </>
          ) : (
            <>
              <div className="w-4 h-4 rounded-md flex items-center justify-center text-2xs"
                style={{ background: 'linear-gradient(135deg, #22c55e, #10b981)' }}>
                🧠
              </div>
              <span className="text-2xs text-text-muted uppercase tracking-wider">U-Hermes</span>
            </>
          )}
        </div>

        <div
          className={isUser
            ? 'px-4 py-2.5 rounded-2xl rounded-br-md text-sm leading-relaxed text-white'
            : 'px-4 py-3.5 rounded-xl rounded-bl-md text-sm leading-relaxed border border-border bg-surface-card'
          }
          style={isUser ? {
            background: 'linear-gradient(135deg, #6366f1, #7c3aed)',
            boxShadow: '0 2px 12px rgba(99,102,241,0.25)',
          } : undefined}
        >
          {isUser ? (
            <p className="whitespace-pre-wrap">{message.content}</p>
          ) : (
            <div className="prose prose-invert prose-sm max-w-none">
              <ReactMarkdown>{message.content}</ReactMarkdown>
            </div>
          )}

          {isStreaming && !isUser && (
            <span className="inline-block w-2 h-4 ml-0.5 align-text-bottom bg-primary animate-pulse rounded-sm" />
          )}
        </div>

        {!isUser && !isStreaming && message.content && (
          <div className="flex gap-3 mt-1.5 px-1">
            <button
              onClick={() => navigator.clipboard.writeText(message.content)}
              className="text-2xs text-text-muted hover:text-text-secondary transition-colors"
            >
              📋 复制
            </button>
          </div>
        )}
      </div>
    </div>
  );
}
