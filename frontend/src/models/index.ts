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
