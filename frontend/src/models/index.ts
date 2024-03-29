export interface User {
  id: number;
  username: string;
}

export interface NumericRecord {
  id: number;
  timestamp: string;
  value: number;
  comment: string;
  // eslint-disable-next-line no-use-before-define
  hobbit?: Hobbit;
}

export interface AppPassword {
    id: string;
    description: string;
    secret: string;
    last_used_at: string;
}

export interface Hobbit {
  id: number;
  user: User;
  name: string;
  image: string;
  description: string;
  records: NumericRecord[];
  heatmap: NumericRecord[];
}

export enum FeedEventTypus {
  HobbitCreated = 'HobbitCreated',
  RecordCreated = 'RecordCreated',
}

export interface FeedEvent {
  FeedEventTypus: FeedEventTypus;
  CreatedAt: string;
  Payload: Hobbit | NumericRecord;
}
