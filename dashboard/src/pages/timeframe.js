import { List, Edit, Datagrid, TextField, BooleanField, 
    NumberInput, BooleanInput, SimpleForm, TextInput, NumberField, Pagination 
} from 'react-admin';



export const TimeframeList = () => (
    <List pagination={<Pagination rowsPerPageOptions={[100]}/>}>
        <Datagrid rowClick="edit">
            <TextField source="format" label="Timeframe"/>
            <NumberField source="percent_alert"/>
            <BooleanField source="default" />
        </Datagrid>
    </List>
);

export const TimeframeEdit = () => (
    <Edit>
        <SimpleForm>
            <TextInput source="id" disabled/>
            <NumberInput source="percent_alert"/>
            <BooleanInput source="default"/>
        </SimpleForm>
    </Edit>
);