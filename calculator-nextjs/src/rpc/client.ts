import {createClient} from "@connectrpc/connect";
import {createConnectTransport} from "@connectrpc/connect-web";
import {CalculatorService} from "@/rpc/calculator_pb";

const transport = createConnectTransport({
    baseUrl: typeof window === 'undefined' ? '' : 'http://localhost:4090', // 处理SSR情况
    useBinaryFormat: true, // 使用二进制格式
});
export const client = createClient(CalculatorService, transport);
