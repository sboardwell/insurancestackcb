import { useState } from 'react';
import { AlertCircle, X, Info, AlertTriangle, CheckCircle } from 'lucide-react';
import useRoxFlag from '../hooks/useRoxFlag';

interface AlertBannerProps {
  type?: 'info' | 'success' | 'warning' | 'critical';
  title: string;
  message: string;
  dismissible?: boolean;
}

export default function AlertBanner({
  type = 'info',
  title,
  message,
  dismissible = true,
}: AlertBannerProps) {
  const [dismissed, setDismissed] = useState(false);
  const alertsBanner = useRoxFlag('alertsBanner');

  // Don't render if feature flag is disabled or banner is dismissed
  if (!alertsBanner || dismissed) {
    return null;
  }

  const styles = {
    info: {
      bg: 'bg-blue-50',
      border: 'border-blue-200',
      text: 'text-blue-800',
      icon: Info,
    },
    success: {
      bg: 'bg-green-50',
      border: 'border-green-200',
      text: 'text-green-800',
      icon: CheckCircle,
    },
    warning: {
      bg: 'bg-yellow-50',
      border: 'border-yellow-200',
      text: 'text-yellow-800',
      icon: AlertTriangle,
    },
    critical: {
      bg: 'bg-red-50',
      border: 'border-red-200',
      text: 'text-red-800',
      icon: AlertCircle,
    },
  };

  const style = styles[type];
  const Icon = style.icon;

  return (
    <div className={`${style.bg} ${style.border} border rounded-lg p-4 mb-6 animate-slide-in`}>
      <div className="flex items-start">
        <Icon className={`w-5 h-5 ${style.text} mt-0.5 flex-shrink-0`} />
        <div className="ml-3 flex-1">
          <h3 className={`text-sm font-semibold ${style.text}`}>{title}</h3>
          <p className={`mt-1 text-sm ${style.text} opacity-90`}>{message}</p>
        </div>
        {dismissible && (
          <button
            onClick={() => setDismissed(true)}
            className={`ml-3 ${style.text} hover:opacity-70 transition-opacity flex-shrink-0`}
            aria-label="Dismiss alert"
          >
            <X className="w-5 h-5" />
          </button>
        )}
      </div>
    </div>
  );
}
