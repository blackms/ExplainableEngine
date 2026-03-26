'use client';

interface ConfidenceGaugeProps {
  confidence: number;
  size?: number;
}

function getConfidenceColor(confidence: number): string {
  if (confidence >= 0.8) return 'text-green-500';
  if (confidence >= 0.5) return 'text-yellow-500';
  return 'text-red-500';
}

function getConfidenceStroke(confidence: number): string {
  if (confidence >= 0.8) return 'stroke-green-500';
  if (confidence >= 0.5) return 'stroke-yellow-500';
  return 'stroke-red-500';
}

function getConfidenceTrack(confidence: number): string {
  if (confidence >= 0.8) return 'stroke-green-500/20';
  if (confidence >= 0.5) return 'stroke-yellow-500/20';
  return 'stroke-red-500/20';
}

export function ConfidenceGauge({ confidence, size = 80 }: ConfidenceGaugeProps) {
  const radius = 34;
  const circumference = 2 * Math.PI * radius;
  const filled = circumference * confidence;
  const percentage = (confidence * 100).toFixed(1);

  return (
    <div className="relative inline-flex items-center justify-center" style={{ width: size, height: size }}>
      <svg
        viewBox="0 0 80 80"
        width={size}
        height={size}
        className="-rotate-90"
      >
        <circle
          cx="40"
          cy="40"
          r={radius}
          fill="none"
          strokeWidth="6"
          className={getConfidenceTrack(confidence)}
        />
        <circle
          cx="40"
          cy="40"
          r={radius}
          fill="none"
          strokeWidth="6"
          strokeLinecap="round"
          strokeDasharray={circumference}
          strokeDashoffset={circumference - filled}
          className={`${getConfidenceStroke(confidence)} transition-all duration-500`}
        />
      </svg>
      <span className={`absolute text-sm font-semibold ${getConfidenceColor(confidence)}`}>
        {percentage}%
      </span>
    </div>
  );
}
