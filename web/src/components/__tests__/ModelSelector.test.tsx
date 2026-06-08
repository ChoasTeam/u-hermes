import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import ModelSelector from '../ModelSelector';

describe('ModelSelector', () => {
  const onSelect = vi.fn();

  const renderSelector = (selectedId = 'deepseek') =>
    render(<ModelSelector selectedId={selectedId} onSelect={onSelect} />);

  it('renders all 4 preset models', () => {
    renderSelector();
    expect(screen.getByText('DeepSeek V3')).toBeInTheDocument();
    expect(screen.getByText('Kimi K2.5')).toBeInTheDocument();
    expect(screen.getByText('通义千问 Qwen')).toBeInTheDocument();
    expect(screen.getByText('智谱 GLM')).toBeInTheDocument();
  });

  it('shows checkmark on selected model', () => {
    const { container } = renderSelector('deepseek');

    const selectedCard = container.querySelector('[class*="border-primary"]');
    expect(selectedCard).toBeTruthy();
  });

  it('calls onSelect when clicking a model', () => {
    renderSelector('deepseek');

    screen.getByText('Kimi K2.5').click();
    expect(onSelect).toHaveBeenCalledWith('kimi');
  });
});
