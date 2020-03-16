export const BASE_URL = "/api/v0";

enum Method {
  GET = "GET"
}

export async function Get<T>(url: string): Promise<T> {
  const response: T = await issueFetch(url, Method.GET);
  return response;
}

async function issueFetch<T>(
  url: string,
  method: Method,
  data?: any
): Promise<T> {
  const response = await fetch(url, {
    body: data ? JSON.stringify(data) : undefined,
    cache: "no-cache",
    headers: {
      "Content-Type": "application/json"
    },
    method,
    referrer: "no-referrer"
  });

  if (!response.ok) {
    const errText = await response.text();
    throw Error(errText ? errText : response.statusText);
  }

  try {
    const responseBody: T = await response.json();
    return responseBody;
  } catch (err) {
    return {} as T; // response was text or empty
  }
}
