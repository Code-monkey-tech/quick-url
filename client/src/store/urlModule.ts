import { requestShortUrl } from "../api/index";
import { IUrlData } from "../components/types";
import { useToast } from "vue-toastification";
import { ActionContext } from "vuex";
import { IRootState } from "../components/types";

export default {
  state: {
    shortUrl: "",
  },
  mutations: {
    setLongUrl(state: { shortUrl: string }, payload: string) {
      state.shortUrl = payload;
    },
  },
  getters: {
    getCorrectShortUrl(state: { shortUrl: string }) {
      const currentUrl = window.location.origin;
      const resultUrl = state.shortUrl ? `${currentUrl}/${state.shortUrl}` : "";
      return resultUrl;
    },
  },
  actions: {
    async getShortUrl(
      { commit }: ActionContext<IRootState, IRootState>,
      requestData: IUrlData
    ) {
      const toast = useToast();
      try {
        const res = await requestShortUrl(requestData);
        if (res.status >= 200 && res.status <= 299) {
          const { data } = res;
          toast.success("Good Request!");
          commit("setLongUrl", data.url);
        }
      } catch (e: any) {
        toast.error(e.message);
      }
    },
  },
};
