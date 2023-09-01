export type Error = { message?: string };
export type User = { name?: string };

export interface GetHelloParams {
  q1?: string;
  q2: string;
  h1?: string;
}
export interface PostHelloParams {}
export interface GetUserParams {}

export class Api {
  constructor(private baseUrl: string = window.location.origin) {}

  async getHello(params: GetHelloParams): Promise<unknown> {
    const headers = new Headers();
    const url = new URL("/", this.baseUrl);
    if (params.q1 !== undefined) {
      url.searchParams.set("q1", params.q1);
    }
    if (params.q2 !== undefined) {
      url.searchParams.set("q2", params.q2);
    }
    if (params.h1 !== undefined) {
      headers.set("h1", params.h1);
    }
    const response = await fetch(url, { headers });
    const body = await response.json();
    return body;
  }
  async postHello(params: PostHelloParams): Promise<unknown> {
    const headers = new Headers();
    const url = new URL("/", this.baseUrl);

    const response = await fetch(url, { headers });
    const body = await response.json();
    return body;
  }
  async getUser(params: GetUserParams): Promise<User> {
    const headers = new Headers();
    const url = new URL("/users", this.baseUrl);

    const response = await fetch(url, { headers });
    const body = await response.json();
    return body;
  }
}
