import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import StepIndicator from '../components/StepIndicator';
import ModelSelector from '../components/ModelSelector';
import { updateModel, testModel } from '../api/client';

const PRESET_MODELS: Record<string, { name: string; apiBase: string }> = {
  deepseek: { name: 'DeepSeek V3', apiBase: 'https://api.deepseek.com' },
  kimi: { name: 'Kimi K2.5', apiBase: 'https://api.moonshot.cn' },
  qwen: { name: '通义千问 Qwen', apiBase: 'https://dashscope.aliyuncs.com/compatible-mode' },
  glm: { name: '智谱 GLM', apiBase: 'https://open.bigmodel.cn/api/paas/v4' },
};

export default function OnboardingPage() {
  const [step, setStep] = useState(0);
  const [selectedModel, setSelectedModel] = useState('deepseek');
  const [apiKey, setApiKey] = useState('');
  const [testing, setTesting] = useState(false);
  const [testResult, setTestResult] = useState<'idle' | 'success' | 'error'>('idle');
  const [testMessage, setTestMessage] = useState('');
  const navigate = useNavigate();

  const handleTest = async () => {
    setTesting(true);
    setTestResult('idle');

    const preset = PRESET_MODELS[selectedModel];
    await updateModel('default', {
      name: preset.name,
      api_base: preset.apiBase,
      api_key: apiKey,
      is_default: true,
    } as any);

    try {
      const result = await testModel('default');
      if (result.status === 'ok') {
        setTestResult('success');
        setTimeout(() => setStep(2), 800);
      } else {
        setTestResult('error');
        setTestMessage(result.message || '连接失败');
      }
    } catch {
      setTestResult('error');
      setTestMessage('网络错误');
    } finally {
      setTesting(false);
    }
  };

  return (
    <div className="min-h-screen bg-surface flex items-center justify-center p-6">
      <div className="w-full max-w-md">
        {step === 0 && (
          <div className="text-center">
            <div className="text-6xl mb-5">🧠</div>
            <h1 className="text-2xl font-bold mb-2">欢迎使用 U-Hermes</h1>
            <p className="text-text-secondary text-sm mb-8 leading-relaxed">
              你的随身 AI 助手<br />插上 U 盘，随时随地使用
            </p>

            <div className="flex gap-3 justify-center mb-8">
              {[
                { icon: '🔌', label: '即插即用' },
                { icon: '🔒', label: '数据本地' },
                { icon: '🧠', label: '越用越聪明' },
              ].map((item) => (
                <div key={item.label} className="bg-surface-card rounded-xl p-4 text-center min-w-[90px]">
                  <div className="text-xl mb-1">{item.icon}</div>
                  <div className="text-2xs text-text-muted">{item.label}</div>
                </div>
              ))}
            </div>

            <button
              onClick={() => setStep(1)}
              className="px-8 py-3 rounded-xl text-sm font-semibold text-white transition-all hover:opacity-90"
              style={{ background: 'linear-gradient(135deg, #6366f1, #7c3aed)' }}
            >
              开始配置 →
            </button>

            <StepIndicator steps={3} current={0} />
          </div>
        )}

        {step === 1 && (
          <div>
            <h2 className="text-lg font-semibold mb-1">选择 AI 模型</h2>
            <p className="text-sm text-text-muted mb-5">推荐使用免费额度大的模型</p>

            <ModelSelector selectedId={selectedModel} onSelect={setSelectedModel} />

            <div className="mt-5">
              <label className="text-2xs text-text-secondary mb-1 block">API Key</label>
              <div className="flex gap-2">
                <input
                  type="password"
                  value={apiKey}
                  onChange={(e) => setApiKey(e.target.value)}
                  placeholder="sk-..."
                  className="flex-1 bg-surface-card border border-border rounded-lg px-3.5 py-2.5 text-sm text-text-primary outline-none focus:border-primary/50"
                />
              </div>
              <p className="text-2xs text-text-muted mt-1.5">
                不知道在哪获取？<a href="#" className="text-primary hover:underline">查看图文教程 →</a>
              </p>
            </div>

            {testResult === 'error' && (
              <div className="mt-4 p-3 rounded-lg bg-red-500/5 border border-red-500/20 text-xs text-red-400">
                {testMessage}
              </div>
            )}
            {testResult === 'success' && (
              <div className="mt-4 p-3 rounded-lg bg-green-500/5 border border-green-500/20 text-xs text-green-400">
                ✓ 连接成功
              </div>
            )}

            <div className="flex justify-between items-center mt-6">
              <button onClick={() => setStep(0)} className="text-sm text-text-muted hover:text-text-secondary">← 上一步</button>
              <button
                onClick={handleTest}
                disabled={!apiKey.trim() || testing}
                className="px-6 py-2.5 rounded-lg text-sm font-medium text-white transition-all disabled:opacity-30"
                style={{ background: 'linear-gradient(135deg, #6366f1, #7c3aed)' }}
              >
                {testing ? '测试中...' : '测试连接 →'}
              </button>
            </div>

            <div className="mt-6">
              <StepIndicator steps={3} current={1} />
            </div>
          </div>
        )}

        {step === 2 && (
          <div className="text-center">
            <div className="w-16 h-16 rounded-full flex items-center justify-center text-3xl mx-auto mb-5"
              style={{
                background: 'linear-gradient(135deg, #22c55e, #10b981)',
                boxShadow: '0 0 30px rgba(34,197,94,0.3)',
              }}>
              ✓
            </div>
            <h2 className="text-xl font-bold mb-1.5">配置完成！</h2>
            <p className="text-sm text-text-muted mb-4 leading-relaxed">
              {PRESET_MODELS[selectedModel].name} 已就绪<br />随时可以开始对话
            </p>
            <div className="inline-block bg-surface-card rounded-lg px-4 py-2.5 mb-6">
              <span className="text-2xs text-text-secondary">💡 以后双击 exe 直接进入聊天，无需重新配置</span>
            </div>
            <br />
            <button
              onClick={() => navigate('/chat')}
              className="px-8 py-3 rounded-xl text-sm font-semibold text-white transition-all hover:opacity-90"
              style={{ background: 'linear-gradient(135deg, #6366f1, #7c3aed)' }}
            >
              开始聊天 →
            </button>

            <div className="mt-6">
              <StepIndicator steps={3} current={2} />
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
