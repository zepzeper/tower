import { useParams } from 'react-router-dom';
import FieldMapping from '../../components/pages/FieldMapping';

function ConnectionEdit() {
  const { id } = useParams();

  return (
    <div className="p-4">
      <FieldMapping />
    </div>
  );
}

export default ConnectionEdit;
