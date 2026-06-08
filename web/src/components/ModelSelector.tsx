interface Props {
  selectedId: string;
  onSelect: (id: string) => void;
}

const PRESET_MODELS = [
  { id: 'deepseek', name: 'DeepSeek V3', desc: '编程能力强 · 注册送 500 万 token', apiBase: 'https://api.deepseek.com' },
  { id: 'kimi', name: 'Kimi K2.5', desc: '超长上下文 · 适合处理长文档', apiBase: 'https://api.moonshot.cn' },
  { id: 'qwen', name: '通义千问 Qwen', desc: '阿里云生态 · 免费额度大', apiBase: 'https://dashscope.aliyuncs.com/compatible-mode' },
  { id: 'glm', name: '智谱 GLM', desc: '学术场景友好', apiBase: 'https://open.bigmodel.cn/api/paas/v4' },
];

export default function ModelSelector({ selectedId, onSelect }: Props) {
  return (
    <div className="space-y-2">
      {PRESET_MODELS.map((preset) => (
        <div
          key={preset.id}
          onClick={() => onSelect(preset.id)}
          className={`p-3.5 rounded-xl cursor-pointer transition-all border-2 ${
            selectedId === preset.id
              ? 'border-primary bg-primary/5'
              : 'border-border bg-surface-card hover:border-primary/30'
          }`}
        >
          <div className="flex justify-between items-center">
            <div>
              <div className="text-sm font-semibold text-text-primary">{preset.name}</div>
              <div className="text-2xs text-text-muted mt-0.5">{preset.desc}</div>
            </div>
            {selectedId === preset.id && (
              <div className="w-5 h-5 rounded-full bg-primary flex items-center justify-center text-2xs text-white">✓</div>
            )}
          </div>
        </div>
      ))}
    </div>
  );
}
