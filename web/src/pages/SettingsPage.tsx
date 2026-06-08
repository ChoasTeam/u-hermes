import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import ResetConfirm from '../components/ResetConfirm';
import { getSettings, updateSettings, resetSettings, listModels, testModel } from '../api/client';

export default function SettingsPage() {
  const [systemPrompt, setSystemPrompt] = useState('');
  const [models, setModels] = useState<any[]>([]);
  const [showReset, setShowReset] = useState(false);
  const [testResults, setTestResults] = useState<Record<string, string>>({});
  const navigate = useNavigate();

  useEffect(() => {
    getSettings().then((s) => setSystemPrompt(s.chat?.system_prompt || ''));
    listModels().then(setModels);
  }, []);

  const handleSavePrompt = async () => {
    await updateSettings({ chat: { system_prompt: systemPrompt } });
  };

  const handleTestModel = async (id: string) => {
    setTestResults((prev) => ({ ...prev, [id]: 'testing' }));
    try {
      const result = await testModel(id);
      setTestResults((prev) => ({ ...prev, [id]: result.status === 'ok' ? 'success' : 'error' }));
    } catch {
      setTestResults((prev) => ({ ...prev, [id]: 'error' }));
    }
  };

  const handleReset = async () => {
    await resetSettings();
    window.location.href = '/onboarding';
  };

  return (
    <div className="min-h-screen bg-surface p-6">
      <div className="max-w-2xl mx-auto">
        <div className="flex items-center justify-between mb-6">
          <h1 className="text-xl font-bold">⚙️ 设置</h1>
          <button onClick={() => navigate('/chat')} className="text-sm text-text-muted hover:text-text-secondary">← 返回聊天</button>
        </div>

        <div className="bg-surface-card border border-border rounded-xl p-5 mb-4">
          <h3 className="font-semibold text-sm mb-4">🤖 AI 模型</h3>
          {models.map((model) => (
            <div key={model.id} className="mb-4 last:mb-0">
              <div className="text-xs text-text-muted mb-1">{model.name}</div>
              <div className="flex gap-2 items-center">
                <input
                  type="text"
                  value={model.api_base || ''}
                  readOnly
                  className="flex-1 bg-surface border border-border rounded-lg px-3 py-2 text-sm text-text-secondary"
                />
                <button
                  onClick={() => handleTestModel(model.id)}
                  className="px-4 py-2 rounded-lg text-xs font-medium bg-surface-hover text-text-secondary hover:text-text-primary"
                >
                  {testResults[model.id] === 'testing' ? '测试中...' :
                   testResults[model.id] === 'success' ? '✓ 正常' :
                   testResults[model.id] === 'error' ? '✗ 失败' : '测试连接'}
                </button>
              </div>
            </div>
          ))}
        </div>

        <div className="bg-surface-card border border-border rounded-xl p-5 mb-4">
          <h3 className="font-semibold text-sm mb-3">💬 聊天设置</h3>
          <div className="text-2xs text-text-muted mb-1.5">系统提示词</div>
          <textarea
            value={systemPrompt}
            onChange={(e) => setSystemPrompt(e.target.value)}
            rows={3}
            className="w-full bg-surface border border-border rounded-lg px-3.5 py-2.5 text-sm text-text-primary outline-none focus:border-primary/50 resize-none"
          />
          <button
            onClick={handleSavePrompt}
            className="mt-3 px-4 py-2 rounded-lg text-xs font-medium text-white"
            style={{ background: 'linear-gradient(135deg, #6366f1, #7c3aed)' }}
          >
            保存
          </button>
        </div>

        <div className="bg-red-500/5 border border-red-500/20 rounded-xl p-5">
          <h3 className="font-semibold text-sm text-red-400 mb-1.5">⚠️ 恢复出厂设置</h3>
          <p className="text-xs text-text-muted mb-4">将删除所有配置、对话历史和数据。此操作不可撤销。</p>
          <button
            onClick={() => setShowReset(true)}
            className="px-5 py-2 rounded-lg text-sm font-medium text-white bg-red-500 hover:bg-red-600 transition-colors"
          >
            恢复出厂设置
          </button>
        </div>

        {showReset && (
          <ResetConfirm
            onConfirm={handleReset}
            onCancel={() => setShowReset(false)}
          />
        )}
      </div>
    </div>
  );
}
