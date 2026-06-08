import { useState, useRef, useCallback } from 'react';
import { streamChat, Message } from '../api/client';

export function useChat() {
  const [messages, setMessages] = useState<Message[]>([]);
  const [isStreaming, setIsStreaming] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const abortRef = useRef<AbortController | null>(null);

  const send = useCallback(async (
    conversationId: string,
    content: string,
    model: string,
  ) => {
    setError(null);
    setIsStreaming(true);

    const userMsg: Message = {
      id: crypto.randomUUID(),
      conversation_id: conversationId,
      role: 'user',
      content,
      tokens_used: 0,
      created_at: Date.now() / 1000,
    };
    setMessages((prev) => [...prev, userMsg]);

    const aiMsgId = crypto.randomUUID();
    const aiMsg: Message = {
      id: aiMsgId,
      conversation_id: conversationId,
      role: 'assistant',
      content: '',
      tokens_used: 0,
      created_at: Date.now() / 1000,
    };
    setMessages((prev) => [...prev, aiMsg]);

    abortRef.current = streamChat(
      conversationId,
      content,
      model,
      (token) => {
        setMessages((prev) =>
          prev.map((m) =>
            m.id === aiMsgId ? { ...m, content: m.content + token } : m,
          ),
        );
      },
      (data) => {
        setMessages((prev) =>
          prev.map((m) =>
            m.id === aiMsgId
              ? { ...m, tokens_used: data.total_tokens }
              : m,
          ),
        );
        setIsStreaming(false);
      },
      (err) => {
        setError(err);
        setIsStreaming(false);
      },
    );

    return aiMsgId;
  }, []);

  const stop = useCallback(() => {
    abortRef.current?.abort();
    setIsStreaming(false);
  }, []);

  const clearMessages = useCallback(() => {
    setMessages([]);
    setError(null);
  }, []);

  return { messages, isStreaming, error, send, stop, setMessages, clearMessages };
}
