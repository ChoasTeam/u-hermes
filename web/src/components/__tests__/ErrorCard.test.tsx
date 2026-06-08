import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import ErrorCard from '../ErrorCard';

describe('ErrorCard', () => {
  it('renders icon, title and description', () => {
    render(
      <ErrorCard icon="⚠️" title="连接失败" description="无法连接到 AI 服务" />
    );

    expect(screen.getByText('⚠️')).toBeInTheDocument();
    expect(screen.getByText('连接失败')).toBeInTheDocument();
    expect(screen.getByText('无法连接到 AI 服务')).toBeInTheDocument();
  });

  it('renders action buttons when provided', () => {
    const onRetry = vi.fn();
    const onSettings = vi.fn();

    render(
      <ErrorCard
        icon="⚠️"
        title="连接失败"
        description="请检查网络"
        actions={[
          { label: '重试', onClick: onRetry, primary: true },
          { label: '修改配置', onClick: onSettings },
        ]}
      />
    );

    const retryBtn = screen.getByText('重试');
    const settingsBtn = screen.getByText('修改配置');

    fireEvent.click(retryBtn);
    expect(onRetry).toHaveBeenCalledOnce();

    fireEvent.click(settingsBtn);
    expect(onSettings).toHaveBeenCalledOnce();
  });

  it('renders without action area when no actions provided', () => {
    const { container } = render(
      <ErrorCard icon="ℹ️" title="提示" description="一切正常" />
    );

    const buttons = container.querySelectorAll('button');
    expect(buttons).toHaveLength(0);
  });
});
