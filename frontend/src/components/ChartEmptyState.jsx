import React from "react";
import "../styles/ChartEmptyState.css";

export default function ChartEmptyState({ onRetry, chosenLabel, mode = "daily" }) {
  // mode "day" or "weekly"
  return (
    <div id="chart-empty-state" className="ces-root" role="status" aria-live="polite">
      <svg
        className="ces-illustration"
        width="140"
        height="100"
        viewBox="0 0 140 100"
        aria-hidden="true"
        focusable="false"
      >
        <defs>
          <linearGradient id="ces-g1" x1="0" x2="1">
            <stop offset="0%" stopColor="#E6F7F6" />
            <stop offset="100%" stopColor="#F6FBFB" />
          </linearGradient>
          <linearGradient id="ces-g2" x1="0" x2="1">
            <stop offset="0%" stopColor="#5F9EA0" stopOpacity="0.15" />
            <stop offset="100%" stopColor="#5F9EA0" stopOpacity="0.03" />
          </linearGradient>
        </defs>

        <rect x="0" y="0" width="140" height="100" rx="10" fill="url(#ces-g1)" />

        <g transform="translate(18,14)">
          <path d="M10 60 C30 30, 70 30, 90 60" fill="url(#ces-g2)" />
          <circle cx="30" cy="38" r="6" fill="#5F9EA0" opacity="0.95" />
          <circle cx="58" cy="28" r="4" fill="#80B3B1" opacity="0.95" />
          <rect x="46" y="46" rx="3" width="34" height="6" fill="#dfeff0" />
        </g>

        <g transform="translate(100,12)">
          <g transform="scale(0.7)">
            <path d="M8 8 L22 8 L22 22 L8 22 Z" fill="#fff" stroke="#e6f2f1" />
            <path d="M11 14 L19 14" stroke="#c6e7e6" strokeWidth="1.5" strokeLinecap="round" />
          </g>
        </g>
      </svg>

      <div className="ces-text">
        <h2 className="ces-title">
          {mode === "weekly" ? "No data found for the chosen week" : "No data found for chosen day"}
        </h2>

        <p className="ces-body">
          {chosenLabel ? (
            <>
              We couldn't find any measurements for <strong className="ces-highlight">{chosenLabel}</strong>.
            </>
          ) : (
            mode === "weekly"
              ? "We couldn't find any measurements for the selected week."
              : "We couldn't find any measurements for the selected day."
          )}
        </p>

        <div className="ces-actions" role="group" aria-label="Empty state actions">
          <button
            type="button"
            onClick={onRetry}
            className="ces-btn ces-btn-primary"
            aria-label="Retry loading data"
          >
            Retry
          </button>

          <button
            type="button"
            onClick={() => window.location.reload()}
            className="ces-btn ces-btn-ghost"
            aria-label="Refresh page"
          >
            Refresh
          </button>
        </div>

        <small className="ces-hint">
          Tip: try a different {mode === "weekly" ? "week" : "date"} or check that your device was connected that {mode === "weekly" ? "week" : "day"}.
        </small>
      </div>
    </div>
  );
}