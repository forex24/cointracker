import * as React from "react";
import { Admin, Resource } from 'react-admin';
import dataProvider from './dataProvider/rest';

import { Klines } from "./pages/kline";
import { AlertConfigList, AlertConfigCreate, AlertConfigEdit,} from "./pages/alertConfig";
import { AlertList } from "./pages/alert";
import { TimeframeList, TimeframeEdit } from "./pages/timeframe";


import AddAlertIcon from '@mui/icons-material/AddAlert';
import CandlestickChartIcon from '@mui/icons-material/CandlestickChart';
import CurrencyBitcoinIcon from '@mui/icons-material/CurrencyBitcoin';
import PermDataSettingIcon from '@mui/icons-material/PermDataSetting';

const App = () => 
<Admin title="CRYPTO TRAKING" disableTelemetry dataProvider={dataProvider('/api/v1')}>
  <Resource name="klines" list={Klines} icon={CandlestickChartIcon}/>
  <Resource name="alerts" list={AlertList} icon={AddAlertIcon}/>
  <Resource name="configs" list={AlertConfigList} create={AlertConfigCreate} edit={AlertConfigEdit} icon={CurrencyBitcoinIcon} options={{ label: 'Coinlists' }}  />
  <Resource name="timeframes" list={TimeframeList} edit={TimeframeEdit} icon={PermDataSettingIcon} options={{ label: 'Configs' }} />
</Admin>;

export default App;