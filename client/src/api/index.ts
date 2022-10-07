import axios from "axios";
import { IUrlData } from "../components/types";
export const API_ROOT = "https://e7ast1c-shrty.herokuapp.com";

export const requestShortUrl = async (requestData: IUrlData): Promise<any> => {
  const url = `${API_ROOT}/shorten`;
  const res = await axios.post<IUrlData>(url, requestData);
  return res;
};

export const requestExpandUrl = async (hash: string): Promise<any> => {
  const url = `${API_ROOT}/expand`;
  const res = await axios.get<IUrlData>(url, {
    params: {
      hash,
    },
  });
  return res.data;
};
