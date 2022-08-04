import { List, Datagrid, TextField, Create, SimpleForm, ReferenceInput, AutocompleteInput, required, 
    Edit, SelectInput,BooleanInput, ReferenceArrayField, DateField, SearchInput,
    SingleFieldList, ChipField , ReferenceArrayInput, Pagination,
    SelectArrayInput, FileInput ,FileField
} from 'react-admin';


const alertConfigsFilters = [
    <SearchInput source="q" alwaysOn />
];


export const AlertConfigList = () => (
    <List filters={alertConfigsFilters} pagination={<Pagination rowsPerPageOptions={[100, 10]} />}>
        <Datagrid rowClick="edit" >
            <TextField source="id" />
            <TextField source="symbol" />
            <TextField source="direction" />
            <ReferenceArrayField label="Timeframes" reference="timeframes" source="timeframes" defaultValue={[5]}>
                <SingleFieldList>
                    <ChipField source="format" />
                </SingleFieldList>
            </ReferenceArrayField>
            <ReferenceArrayField label="Disabled" reference="timeframes" source="disabled_timeframes">
                <SingleFieldList>
                    <ChipField source="format" />
                </SingleFieldList>
            </ReferenceArrayField>
            <DateField source="last_triggered_at" showTime/>
        </Datagrid>
    </List>
);

export const AlertConfigCreate = () => (
    <Create redirect="list">
        <SimpleForm>
            <ReferenceInput label="Symbol" source="symbol" reference="symbols"  validate={[required()]}>
                <AutocompleteInput fullWidth optionText="format" optionValue="format"/>
            </ReferenceInput>
            <ReferenceArrayInput label="Timeframes" source="timeframes" reference="timeframes" validate={[required()]}>
                <SelectArrayInput optionText="format" />
            </ReferenceArrayInput>
            <SelectInput source="direction" validate={required()} choices={[
                { id: 'up', name: 'Up' },
                { id: 'down', name: 'Down' },
                { id: 'both', name: 'Both' },
            ]} defaultValue={"up"} />
            <BooleanInput source="auto_disable_after_trigger" />
            <FileInput source="file" label="Import" accept="text/plain" multiple={false}>
                <FileField source="src" title="title" />
            </FileInput>
        </SimpleForm>
    </Create>
);


export const AlertConfigEdit = () => (
    <Edit mutationMode="pessimistic">
        <SimpleForm>
            <ReferenceInput label="Symbol" source="symbol" reference="symbols"  validate={[required()]}>
                <AutocompleteInput fullWidth optionText="format" optionValue="format"/>
            </ReferenceInput>
            <ReferenceArrayInput label="Timeframes" source="timeframes" reference="timeframes" validate={[required()]}>
                <SelectArrayInput optionText="format" />
            </ReferenceArrayInput>
            <ReferenceArrayInput label="DisabledTimeframes" source="disabled_timeframes" reference="timeframes" validate={[required()]}>
                <SelectArrayInput optionText="format" />
            </ReferenceArrayInput>
            <SelectInput source="direction" validate={required()} choices={[
                { id: 'up', name: 'Up' },
                { id: 'down', name: 'Down' },
                { id: 'both', name: 'Both' },
            ]}/>
            <BooleanInput source="auto_disable_after_trigger" />
        </SimpleForm>
    </Edit>
);