import { useState, useEffect, useCallback } from 'react';
import { listConversations, createConversation, deleteConversation, Conversation } from '../api/client';

export function useConversations() {
  const [conversations, setConversations] = useState<Conversation[]>([]);
  const [loading, setLoading] = useState(true);

  const refresh = useCallback(async () => {
    try {
      const list = await listConversations();
      setConversations(list);
    } catch (err) {
      console.error('Failed to list conversations:', err);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => { refresh(); }, [refresh]);

  const create = async () => {
    const conv = await createConversation();
    setConversations((prev) => [conv, ...prev]);
    return conv;
  };

  const remove = async (id: string) => {
    await deleteConversation(id);
    setConversations((prev) => prev.filter((c) => c.id !== id));
  };

  return { conversations, loading, refresh, create, remove };
}
