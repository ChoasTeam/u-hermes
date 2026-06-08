import { useEffect, useRef } from 'react';
import { Message } from '../api/client';
import MessageBubble from './MessageBubble';
import EmptyState from './EmptyState';

interface Props {
  messages: Message[];
  isStreaming: boolean;
}

export default function ChatMessages({ messages, isStreaming }: Props) {
  const bottomRef = useRef<HTMLDivElement>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const userScrolledUp = useRef(false);

  useEffect(() => {
    if (!userScrolledUp.current) {
      bottomRef.current?.scrollIntoView({ behavior: 'smooth' });
    }
  }, [messages]);

  const handleScroll = () => {
    const el = containerRef.current;
    if (!el) return;
    const atBottom = el.scrollHeight - el.scrollTop - el.clientHeight < 60;
    userScrolledUp.current = !atBottom;
  };

  if (messages.length === 0) {
    return <EmptyState />;
  }

  return (
    <div ref={containerRef} onScroll={handleScroll} className="flex-1 overflow-y-auto px-6 py-5">
      {messages.map((msg, i) => (
        <MessageBubble
          key={msg.id}
          message={msg}
          isStreaming={isStreaming && i === messages.length - 1 && msg.role === 'assistant'}
        />
      ))}
      <div ref={bottomRef} />
    </div>
  );
}
