import { useMemo } from 'react';

interface Props {
  onSuggestionClick?: (text: string) => void;
}

const SUGGESTIONS: Record<string, string[]> = {
  morning: ['📝 帮我写一份今日计划', '💡 头脑风暴一个新想法', '📖 总结昨天的会议纪要', '🔧 帮我调试一段代码', '🌐 翻译一份文档', '📊 分析数据趋势'],
  afternoon: ['📝 帮我写日报总结', '📖 总结这篇文章的要点', '🔧 解释这段代码的逻辑', '💡 帮我优化工作方案', '📊 生成数据可视化建议', '✉️ 帮我起草一封邮件'],
  evening: ['📝 帮我复盘今天的工作', '💡 规划明天的任务', '📖 学习一个新概念', '🔧 写一个自动化脚本', '📊 整理本周数据', '🧠 帮我分析一个问题'],
};

export default function EmptyState({ onSuggestionClick }: Props) {
  const suggestions = useMemo(() => {
    const hour = new Date().getHours();
    if (hour < 12) return SUGGESTIONS.morning;
    if (hour < 18) return SUGGESTIONS.afternoon;
    return SUGGESTIONS.evening;
  }, []);

  return (
    <div className="flex-1 flex flex-col items-center justify-center px-6 py-8">
      <div className="text-5xl mb-4">🧠</div>
      <h2 className="text-lg font-semibold mb-1">有什么可以帮你的？</h2>
      <p className="text-sm text-text-muted mb-7">AI 已就绪，试着问点什么</p>

      <div className="flex flex-wrap gap-2 justify-center max-w-lg mb-8">
        {suggestions.map((text) => (
          <button
            key={text}
            onClick={() => onSuggestionClick?.(text)}
            className="px-4 py-2.5 rounded-full text-xs text-text-secondary bg-surface-card border border-border hover:border-primary/30 hover:text-text-primary transition-all"
          >
            {text}
          </button>
        ))}
      </div>
    </div>
  );
}
