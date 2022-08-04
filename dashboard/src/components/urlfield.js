import * as React from 'react';
import Link from '@mui/material/Link';
import { useRecordContext } from 'react-admin';

var currencies = ["USDT", "USD", "BTC", "ETH", "BUSD", "USDC", "DAI", "TUSD", "USDN", "LUSD", "USDD", "USDP"]

function separatePair(pair) {
    for (const currency of currencies) {
        if (pair.endsWith(currency)) {
            var partern = currency + "$"
            var re = new RegExp(partern, "g");
            return {
                coin: pair.replace(re, ""),
                currency: currency,
            }
        }
    }
    return {
        coin: pair,
        currency: "",
    }
}



function getExchangeURL(symbol) {
    const sarr = symbol.split(':')
    if (sarr.length !== 2 ){
        return ""
    }
    const exchange = sarr[0]
    const pair = sarr[1]
    const { coin, currency }  = separatePair(pair)

    switch(exchange){
        case "BINANCE":
            return "https://www.binance.com/en/trade/" + coin + "_" + currency
        case "GATEIO":
            return "https://www.gate.io/trade/" + coin + "_" + currency
        case "COINBASE":
            return "https://pro.coinbase.com/trade/" + coin + "_" + currency
        case "KUCOIN":
            return "https://www.kucoin.com/vi/trade/" + coin + "-" + currency
        case "BITFINEX":
            return "https://trading.bitfinex.com/t/" + coin + "_" + currency
        case "HUOBI":
            return "https://www.huobi.com/en-us/exchange/" + coin + "_" + currency
        case "OKX":
        case "OKEX":
            return "https://www.okx.com/vi/trade-spot/" + coin + "-" + currency
        case "FTX":
            return "https://ftx.com/trade/" + coin + "/" + currency
        case "MEXC":
            return "https://www.mexc.com/vi-VN/exchange/" + coin + "_" + currency
        default:
            return "https://www.tradingview.com/chart?symbol=" + symbol
    }
}

const ExUrlField = ( props ) =>
{
    const { source, target, rel } = props;
    const record = useRecordContext(props);
    const symbol = record && record[source];
    if (symbol == null) {
        return null;
    }
    if (typeof symbol !== 'string' && symbol instanceof String) {   
        return null
    }
    const exchangeTradeURL = getExchangeURL(symbol)
    return (
        <Link href={exchangeTradeURL} target={target} rel={rel}>
            {symbol}
        </Link>
    );
};

ExUrlField.defaultProps = {
    addLabel: true,
};

export default ExUrlField;
