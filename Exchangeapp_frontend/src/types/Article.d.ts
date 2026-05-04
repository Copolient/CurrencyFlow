export interface Article {
  ID: string;
  Title: string;
  Preview: string;
  Content: string;
  CreatedAt?: string;
}

export interface LikeResponse {
  likes: string; // Redis INCR 返回的是字符串
  message?: string;
}

export interface ExchangeRate {
  _id?: number;
  fromCurrency: string;
  toCurrency: string;
  rate: number;
  date?: string;
}

export interface ApiError {
  error: string;
}
