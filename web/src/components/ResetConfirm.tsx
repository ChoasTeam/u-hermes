interface Props {
  onConfirm: () => void;
  onCancel: () => void;
}

export default function ResetConfirm({ onConfirm, onCancel }: Props) {
  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-6">
      <div className="bg-surface-card border border-border rounded-2xl p-6 max-w-sm w-full">
        <h3 className="text-lg font-semibold text-text-primary mb-2">确认恢复出厂设置？</h3>
        <p className="text-sm text-text-secondary mb-6">
          此操作不可撤销，将删除所有配置、对话历史和数据。U-Hermes 将重启并进入初始配置流程。
        </p>
        <div className="flex gap-3 justify-end">
          <button
            onClick={onCancel}
            className="px-5 py-2.5 rounded-lg text-sm bg-surface-hover text-text-secondary hover:text-text-primary transition-colors"
          >
            取消
          </button>
          <button
            onClick={onConfirm}
            className="px-5 py-2.5 rounded-lg text-sm font-medium text-white bg-red-500 hover:bg-red-600 transition-colors"
          >
            确认重置
          </button>
        </div>
      </div>
    </div>
  );
}
