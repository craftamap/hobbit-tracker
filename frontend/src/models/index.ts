export interface Hobbit {
  id: number;
  user: User;
  name: string;
  image: string;
  description: string;
  records: object[];
};

export interface User {
  id: number;
  name: string;
};

export interface NumericRecord {
  id: number;
  timestamp: string;
  value: number;
  comment: string;
}
