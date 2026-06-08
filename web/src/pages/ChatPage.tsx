import { useState, useCallback, useEffect } from 'react';
import TopBar from '../components/TopBar';
import SidePanel from '../components/SidePanel';
import ChatMessages from '../components/ChatMessages';
import ChatInput from '../components/ChatInput';
import ErrorCard from '../components/ErrorCard';
import { useChat } from '../hooks/useChat';
import { useConversations } from '../hooks/useConversations';
import { getConversation, health } from '../api/client';

export default function ChatPage() {
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const [currentConvId, setCurrentConvId] = useState<string | null>(null);
  const [modelName, setModelName] = useState('');
  const [modelConnected, setModelConnected] = useState(false);
  const [needsConfig, setNeedsConfig] = useState(false);

  const { messages, isStreaming, error, send, stop, setMessages, clearMessages } = useChat();
  const { conversations, refresh: refreshConvs, remove } = useConversations();

  useEffect(() => {
    health().then((h) => {
      setModelName(h.model || '');
      setModelConnected(h.model_connected);
      setNeedsConfig(!h.configured);
      if (!h.configured) {
        window.location.href = '/onboarding';
      }
    });
  }, []);

  const loadConversation = useCallback(async (id: string) => {
    setCurrentConvId(id);
    try {
      const data = await getConversation(id);
      setMessages(data.messages || []);
    } catch (err) {
      console.error('Failed to load conversation:', err);
    }
  }, [setMessages]);

  const handleNewConversation = useCallback(async () => {
    clearMessages();
    setCurrentConvId(null);
    setSidebarOpen(false);
  }, [clearMessages]);

  const handleSelectConversation = useCallback((id: string) => {
    loadConversation(id);
  }, [loadConversation]);

  const handleDeleteConversation = useCallback(async (id: string) => {
    await remove(id);
    if (currentConvId === id) {
      clearMessages();
      setCurrentConvId(null);
    }
    refreshConvs();
  }, [remove, currentConvId, clearMessages, refreshConvs]);

  const handleSend = useCallback(async (content: string) => {
    const convId = currentConvId || '';
    await send(convId, content, '', (newConvId: string) => {
      if (!currentConvId) {
        setCurrentConvId(newConvId);
        refreshConvs();
      }
    });
  }, [send, currentConvId, refreshConvs]);

  return (
    <div className="h-screen flex flex-col bg-surface">
      <TopBar
        modelName={modelName}
        modelConnected={modelConnected}
        onToggleSidebar={() => setSidebarOpen(!sidebarOpen)}
      />

      <SidePanel
        isOpen={sidebarOpen}
        conversations={conversations}
        currentId={currentConvId}
        onSelectConversation={handleSelectConversation}
        onNewConversation={handleNewConversation}
        onDeleteConversation={handleDeleteConversation}
        onClose={() => setSidebarOpen(false)}
      />

      <div className="flex-1 flex flex-col min-h-0">
        {error && (
          <div className="px-5 pt-4">
            <ErrorCard
              icon="⚠️"
              title="请求出错"
              description={error}
              actions={[
                { label: '重试', onClick: () => {/* handled by retry */}, primary: true },
              ]}
            />
          </div>
        )}

        <ChatMessages messages={messages} isStreaming={isStreaming} />

        <ChatInput
          onSend={handleSend}
          onStop={stop}
          isStreaming={isStreaming}
          disabled={needsConfig}
        />
      </div>
    </div>
  );
}
