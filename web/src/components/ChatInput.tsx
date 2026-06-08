import { useState, useRef, useEffect } from 'react';

interface Props {
  onSend: (message: string) => void;
  onStop: () => void;
  isStreaming: boolean;
  disabled: boolean;
}

export default function ChatInput({ onSend, onStop, isStreaming, disabled }: Props) {
  const [input, setInput] = useState('');
  const textareaRef = useRef<HTMLTextAreaElement>(null);

  useEffect(() => {
    if (!isStreaming) {
      textareaRef.current?.focus();
    }
  }, [isStreaming]);

  const handleSubmit = () => {
    const trimmed = input.trim();
    if (!trimmed || isStreaming || disabled) return;
    onSend(trimmed);
    setInput('');
    if (textareaRef.current) {
      textareaRef.current.style.height = 'auto';
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSubmit();
    }
  };

  const handleInput = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    setInput(e.target.value);
    const el = e.target;
    el.style.height = 'auto';
    el.style.height = Math.min(el.scrollHeight, 200) + 'px';
  };

  return (
    <div className="border-t border-border px-5 py-3.5" style={{ background: '#0a0a0b' }}>
      <div className="flex items-end gap-2.5 max-w-3xl mx-auto">
        <div className="flex-1 flex items-center bg-surface-card border border-border rounded-xl px-4 py-1.5 focus-within:border-primary/50 transition-colors">
          <textarea
            ref={textareaRef}
            value={input}
            onChange={handleInput}
            onKeyDown={handleKeyDown}
            placeholder={disabled ? '请先配置 AI 模型' : '给 U-Hermes 发消息...'}
            rows={1}
            disabled={isStreaming || disabled}
            className="flex-1 bg-transparent border-none outline-none text-sm text-text-primary placeholder-text-muted resize-none py-2"
          />
          {isStreaming ? (
            <button
              onClick={onStop}
              className="px-4 py-1.5 rounded-lg text-xs font-medium bg-red-500/10 text-red-400 hover:bg-red-500/20 transition-colors"
            >
              停止
            </button>
          ) : (
            <button
              onClick={handleSubmit}
              disabled={!input.trim() || disabled}
              className="px-4 py-1.5 rounded-lg text-xs font-medium text-white transition-all disabled:opacity-30"
              style={{ background: 'linear-gradient(135deg, #6366f1, #7c3aed)' }}
            >
              发送
            </button>
          )}
        </div>
      </div>
      <div className="flex gap-5 justify-center mt-2 text-2xs text-text-muted">
        <span>Enter 发送</span>
        <span>Shift+Enter 换行</span>
      </div>
    </div>
  );
}
