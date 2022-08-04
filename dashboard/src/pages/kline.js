import { List, Datagrid, TextField, SearchInput, ReferenceField, NumberField, DateField, Pagination } from 'react-admin';
import ExUrlField from '../components/urlfield';
import PercentField from '../components/percentfield';

const klineFilter = [
    <SearchInput source="q" alwaysOn />
];


export const Klines = () => (
    <List  filters={klineFilter} pagination={<Pagination rowsPerPageOptions={[10, 50]} />}>
        <Datagrid>
            <TextField source="id" />
            <ExUrlField source="symbol" target='_blank'/>
            <ReferenceField label="Timeframe" source="timeframe" reference="timeframes"  link={false}>
                <TextField source="format" />
            </ReferenceField>
            <PercentField source="percent_changed"/>
            <NumberField source="open" label='Open price'/>
            <NumberField source="close" label='Close price'/>
            <DateField source="created_at" showTime/>
            <DateField source="openned_at" showTime/>
            <DateField source="closed_at" showTime/>
        </Datagrid>
    </List>
);