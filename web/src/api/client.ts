const BASE = '';

interface Model {
  id: string;
  name: string;
  api_base: string;
  is_default: boolean;
  has_key: boolean;
}

interface Conversation {
  id: string;
  title: string;
  created_at: number;
  updated_at: number;
}

interface Message {
  id: string;
  conversation_id: string;
  role: 'user' | 'assistant';
  content: string;
  tokens_used: number;
  created_at: number;
}

export async function health(): Promise<any> {
  const res = await fetch(`${BASE}/api/health`);
  return res.json();
}

export async function listModels(): Promise<Model[]> {
  const res = await fetch(`${BASE}/api/models`);
  return res.json();
}

export async function updateModel(id: string, data: Partial<Model>): Promise<void> {
  await fetch(`${BASE}/api/models/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  });
}

export async function testModel(id: string): Promise<{ status: string; message?: string }> {
  const res = await fetch(`${BASE}/api/models/${id}/test`, { method: 'POST' });
  return res.json();
}

export async function listConversations(): Promise<Conversation[]> {
  const res = await fetch(`${BASE}/api/conversations`);
  return res.json();
}

export async function createConversation(): Promise<Conversation> {
  const res = await fetch(`${BASE}/api/conversations`, { method: 'POST' });
  return res.json();
}

export async function getConversation(id: string): Promise<{ conversation: Conversation; messages: Message[] }> {
  const res = await fetch(`${BASE}/api/conversations/${id}`);
  return res.json();
}

export async function deleteConversation(id: string): Promise<void> {
  await fetch(`${BASE}/api/conversations/${id}`, { method: 'DELETE' });
}

export async function getSettings(): Promise<any> {
  const res = await fetch(`${BASE}/api/settings`);
  return res.json();
}

export async function updateSettings(data: any): Promise<void> {
  await fetch(`${BASE}/api/settings`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  });
}

export async function resetSettings(): Promise<void> {
  await fetch(`${BASE}/api/settings/reset`, { method: 'POST' });
}

export function streamChat(
  conversationId: string,
  message: string,
  model: string,
  onToken: (token: string) => void,
  onDone: (data: { total_tokens: number; conversation_id: string; message_id: string }) => void,
  onError: (error: string) => void,
): AbortController {
  const controller = new AbortController();

  fetch(`${BASE}/api/chat`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ conversation_id: conversationId, message, model }),
    signal: controller.signal,
  }).then(async (response) => {
    const reader = response.body?.getReader();
    if (!reader) return;

    const decoder = new TextDecoder();
    let buffer = '';

    while (true) {
      const { done, value } = await reader.read();
      if (done) break;

      buffer += decoder.decode(value, { stream: true });
      const lines = buffer.split('\n');
      buffer = lines.pop() || '';

      let currentEvent = '';
      for (const line of lines) {
        if (line.startsWith('event: ')) {
          currentEvent = line.slice(7);
          continue;
        }
        if (line.startsWith('data: ')) {
          try {
            const data = JSON.parse(line.slice(6));
            if (currentEvent === 'token' && data.token) {
              onToken(data.token);
            } else if (currentEvent === 'done') {
              onDone(data);
            } else if (currentEvent === 'error') {
              onError(data.error);
            }
          } catch {}
        }
      }
    }
  }).catch((err) => {
    if (err.name !== 'AbortError') onError(err.message);
  });

  return controller;
}

export type { Model, Conversation, Message };
