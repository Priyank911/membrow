export type SearchItem = {
  title: string;
  url: string;
  snippet: string;
};

export type TraceStep = {
  name: string;
  status: "pending" | "running" | "success" | "failed";
  started_at?: string;
  ended_at?: string;
  error?: string;
};

export type RunRecord = {
  id: string;
  query: string;
  status: "queued" | "running" | "succeeded" | "failed";
  result?: {
    search_results: SearchItem[];
    synthesis: {
      answer: string;
      sources: Array<{
        title: string;
        url: string;
        snippet: string;
        text: string;
      }>;
    };
  };
  trace: TraceStep[];
  error?: {
    code: string;
    message: string;
    retryable: boolean;
  };
};
