import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import StepIndicator from '../StepIndicator';

describe('StepIndicator', () => {
  it('renders correct number of steps', () => {
    render(<StepIndicator steps={3} current={0} />);
    const dots = document.querySelectorAll('[class*="rounded-full"]');
    expect(dots).toHaveLength(3);
  });

  it('highlights current step with wider bar', () => {
    const { container } = render(<StepIndicator steps={3} current={1} />);
    const dots = container.querySelectorAll('[class*="rounded-full"]');

    // Current step (index 1) should be wider (w-6)
    expect(dots[1].className).toContain('w-6');
    // Other steps should be small dots (w-1.5)
    expect(dots[0].className).toContain('w-1.5');
  });

  it('renders all inactive when current is out of range', () => {
    const { container } = render(<StepIndicator steps={2} current={5} />);
    const dots = container.querySelectorAll('[class*="rounded-full"]');
    dots.forEach((dot) => {
      expect(dot.className).not.toContain('w-6');
    });
  });
});
