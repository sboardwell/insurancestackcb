// Type definitions for rox-browser
declare module 'rox-browser' {
  export class Flag {
    constructor(defaultValue: boolean);
    isEnabled(): boolean;
  }

  export interface FetcherResults {
    hasChanges: boolean;
    fetcherStatus: string;
  }

  export interface RoxSetupOptions {
    debugLevel?: string;
    configurationFetchedHandler?: (results: FetcherResults) => void;
  }

  export function register(namespace: string, container: any): void;
  export function setup(apiKey?: string, options?: RoxSetupOptions): Promise<void>;
  export function fetch(): Promise<void>;
  export function unfreeze(namespace?: string): void;

  const Rox: {
    Flag: typeof Flag;
    register: typeof register;
    setup: typeof setup;
    fetch: typeof fetch;
    unfreeze: typeof unfreeze;
    RoxSetupOptions: RoxSetupOptions;
  };

  export default Rox;
}
