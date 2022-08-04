import * as React from 'react';
import { useRecordContext, FunctionField } from 'react-admin';

const PercentField = ( props ) =>
{
    const { source, label } = props;
    const record = useRecordContext(props);
    const percent = record && record[source];
    const absPercent = Math.abs(percent)

    return (
        <FunctionField sx={{ color: percent>0?'#16c784':'#ea3943' }}
          label={label}
          render={ () => {
            return `${absPercent}%`
          }}
        />
    );
};

PercentField.defaultProps = {
    addLabel: true,
};

export default PercentField;
