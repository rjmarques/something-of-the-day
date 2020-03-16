import { BASE_URL, Get } from "../shared/services/http";

export interface ISomething {
  Id: number;
  CreatedAt: string;
  Text: string;
}

export async function GetSomething(): Promise<ISomething> {
  const url = `${BASE_URL}/something`;

  const something: ISomething = await Get(url);
  return something;
}
