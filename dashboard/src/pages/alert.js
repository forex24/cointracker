import { List, Datagrid, TextField, ReferenceField, DateField, SearchInput, Pagination, NumberField } from 'react-admin';
import ExUrlField from '../components/urlfield';
import PercentField from '../components/percentfield';

const alertFilters = [
    <SearchInput source="q" alwaysOn/>
];


export const AlertList = () => (
    <List filters={alertFilters} pagination={<Pagination rowsPerPageOptions={[100, 50]} />} bulkActionButtons={false}>
        <Datagrid>
            <TextField source="id" />
            <ExUrlField source="symbol" target='_blank'/>
            <ReferenceField label="Timeframe" source="timeframe" reference="timeframes"  link={false}>
                <TextField source="format" />
            </ReferenceField>
            <PercentField source="percent_changed" label='Percent changed'/>
            <NumberField source="open" label='Open price'/>
            <NumberField source="close" label='Close price'/>
            <DateField source="created_at" showTime/>
        </Datagrid>
    </List>
);