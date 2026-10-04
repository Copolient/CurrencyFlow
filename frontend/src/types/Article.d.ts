export interface Article {
  id: number;
  title: string;
  preview: string;
  content: string;
  createdAt?: string;
}

export interface LikeResponse {
  likes: number;
  message?: string;
}

export interface ExchangeRate {
  id?: number;
  fromCurrency: string;
  toCurrency: string;
  rate: number;
  date?: string;
}

export interface ApiError {
  error: string;
}
