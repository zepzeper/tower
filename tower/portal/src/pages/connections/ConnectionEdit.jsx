import { useParams } from 'react-router-dom';

function ConnectionEdit() {
  const { id } = useParams();

  return (
    <div className="p-4">
      <h1 className="text-xl font-bold">Integration: {id}</h1>
      {/* Fetch config or render based on ID here */}
    </div>
  );
}

export default ConnectionEdit;
