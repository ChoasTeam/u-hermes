import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import ChatInput from '../ChatInput';

describe('ChatInput', () => {
  const defaultProps = {
    onSend: vi.fn(),
    onStop: vi.fn(),
    isStreaming: false,
    disabled: false,
  };

  const renderInput = (props = {}) =>
    render(<ChatInput {...defaultProps} {...props} />);

  it('shows normal placeholder when enabled', () => {
    renderInput();
    expect(screen.getByPlaceholderText('给 U-Hermes 发消息...')).toBeInTheDocument();
  });

  it('shows config prompt placeholder when disabled', () => {
    renderInput({ disabled: true });
    expect(screen.getByPlaceholderText('请先配置 AI 模型')).toBeInTheDocument();
  });

  it('sends message on Enter key', () => {
    const onSend = vi.fn();
    renderInput({ onSend });

    const textarea = screen.getByPlaceholderText('给 U-Hermes 发消息...');
    fireEvent.change(textarea, { target: { value: '你好' } });
    fireEvent.keyDown(textarea, { key: 'Enter', shiftKey: false });

    expect(onSend).toHaveBeenCalledWith('你好');
  });

  it('does not send on Shift+Enter', () => {
    const onSend = vi.fn();
    renderInput({ onSend });

    const textarea = screen.getByPlaceholderText('给 U-Hermes 发消息...');
    fireEvent.change(textarea, { target: { value: '你好' } });
    fireEvent.keyDown(textarea, { key: 'Enter', shiftKey: true });

    expect(onSend).not.toHaveBeenCalled();
  });

  it('shows stop button when streaming', () => {
    renderInput({ isStreaming: true });
    expect(screen.getByText('停止')).toBeInTheDocument();
    expect(screen.queryByText('发送')).not.toBeInTheDocument();
  });

  it('clears input after send', () => {
    const onSend = vi.fn();
    renderInput({ onSend });

    const textarea = screen.getByPlaceholderText('给 U-Hermes 发消息...') as HTMLTextAreaElement;
    fireEvent.change(textarea, { target: { value: '你好' } });
    fireEvent.keyDown(textarea, { key: 'Enter' });

    expect(textarea.value).toBe('');
  });
});
