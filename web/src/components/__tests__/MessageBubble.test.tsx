import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import MessageBubble from '../MessageBubble';

describe('MessageBubble', () => {
  beforeEach(() => {
    // Mock clipboard API
    Object.assign(navigator, {
      clipboard: { writeText: vi.fn().mockResolvedValue(undefined) },
    });
  });

  it('renders user message with "你" label', () => {
    const msg = { id: '1', conversation_id: 'c1', role: 'user' as const, content: '你好', tokens_used: 0, created_at: 0 };
    render(<MessageBubble message={msg} />);

    expect(screen.getAllByText('你').length).toBeGreaterThanOrEqual(1);
    expect(screen.getByText('你好')).toBeInTheDocument();
  });

  it('renders assistant message with "U-Hermes" label', () => {
    const msg = { id: '2', conversation_id: 'c1', role: 'assistant' as const, content: '你好！有什么可以帮你的？', tokens_used: 0, created_at: 0 };
    render(<MessageBubble message={msg} />);

    expect(screen.getByText('U-Hermes')).toBeInTheDocument();
    expect(screen.getByText('你好！有什么可以帮你的？')).toBeInTheDocument();
  });

  it('shows copy button for finished assistant messages', () => {
    const msg = { id: '2', conversation_id: 'c1', role: 'assistant' as const, content: '回复内容', tokens_used: 0, created_at: 0 };
    render(<MessageBubble message={msg} />);

    expect(screen.getByText('📋 复制')).toBeInTheDocument();
  });

  it('hides copy button while streaming', () => {
    const msg = { id: '2', conversation_id: 'c1', role: 'assistant' as const, content: '正在生成...', tokens_used: 0, created_at: 0 };
    render(<MessageBubble message={msg} isStreaming />);

    expect(screen.queryByText('📋 复制')).not.toBeInTheDocument();
  });

  it('shows blinking cursor while streaming assistant message', () => {
    const msg = { id: '2', conversation_id: 'c1', role: 'assistant' as const, content: '生成中', tokens_used: 0, created_at: 0 };
    const { container } = render(<MessageBubble message={msg} isStreaming />);

    // Streaming cursor: animate-pulse div
    const cursor = container.querySelector('.animate-pulse');
    expect(cursor).toBeTruthy();
  });

  it('renders markdown in assistant messages', () => {
    const msg = { id: '2', conversation_id: 'c1', role: 'assistant' as const, content: '**bold** and `code`', tokens_used: 0, created_at: 0 };
    const { container } = render(<MessageBubble message={msg} />);

    // ReactMarkdown renders <strong> for **bold**
    const strong = container.querySelector('strong');
    expect(strong).toBeTruthy();
    expect(strong!.textContent).toBe('bold');

    // ReactMarkdown renders <code> for `code`
    const code = container.querySelector('code');
    expect(code).toBeTruthy();
    expect(code!.textContent).toBe('code');
  });

  it('does not show label or copy for empty assistant message', () => {
    const msg = { id: '2', conversation_id: 'c1', role: 'assistant' as const, content: '', tokens_used: 0, created_at: 0 };
    render(<MessageBubble message={msg} />);

    // No copy button for empty messages
    expect(screen.queryByText('📋 复制')).not.toBeInTheDocument();
  });
});
