interface Props {
  steps: number;
  current: number;
}

export default function StepIndicator({ steps, current }: Props) {
  return (
    <div className="flex gap-1.5 justify-center">
      {Array.from({ length: steps }).map((_, i) => (
        <div
          key={i}
          className={`rounded-full transition-all ${i === current ? 'w-6 h-1.5 bg-primary' : 'w-1.5 h-1.5 bg-text-muted/30'}`}
        />
      ))}
    </div>
  );
}
