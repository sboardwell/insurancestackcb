import { useState, useEffect } from 'react';
import { getFlagsSnapshot, subscribeFlags } from '../features/flags';

/**
 * Reactive hook for CloudBees Feature Management flags
 * Updates automatically when flag values change in FM dashboard
 *
 * @param key - The flag key to watch
 * @returns Current boolean value of the flag
 *
 * @example
 * const isEnabled = useRoxFlag('claimsFilters');
 */
export default function useRoxFlag(key: string): boolean {
  const [val, setVal] = useState(() => {
    const snapshot = getFlagsSnapshot();
    const initialVal = !!snapshot[key];
    console.log(`[useRoxFlag] Initial value for '${key}':`, initialVal, 'snapshot:', snapshot);
    return initialVal;
  });

  useEffect(() => {
    return subscribeFlags((reason, snap) => {
      const newVal = !!snap[key];
      console.log(`[useRoxFlag] Flag '${key}' update (${reason}):`, newVal, 'snapshot:', snap);
      // Only update if the value actually changed to prevent unnecessary re-renders
      setVal((prevVal) => {
        if (prevVal !== newVal) {
          console.log(`[useRoxFlag] Flag '${key}' changed from ${prevVal} to ${newVal}`);
          return newVal;
        }
        return prevVal;
      });
    });
  }, [key]);

  return val;
}
