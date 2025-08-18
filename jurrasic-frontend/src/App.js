import React, { useState, useEffect } from 'react';
import axios from 'axios';
import 'bootstrap/dist/css/bootstrap.min.css';

const API_BASE_URL = '/dino';

function App() {
  const [dinosaurus, setDinosaurus] = useState([]);
  const [newDino, setNewDino] = useState({
    species: '',
    types: '',
    height: 0,
    length: 0,
    weight: 0,
    aquatic: false,
    flying: false,
  });
  const [editingDino, setEditingDino] = useState(null);

  // Загрузка списка (READ)
  useEffect(() => {
    axios.get(`${API_BASE_URL}`)
      .then(response => setDinosaurus(response.data))
      .catch(error => console.error('Error fetching dinosaurus:', error));
  }, []);

  // Создание (CREATE)
  const handleCreate = () => {
    axios.post(`${API_BASE_URL}`, newDino)
      .then(response => {
        setDinosaurus([...dinosaurus, response.data]);
        setNewDino({
          species: '',
          types: '',
          height: 0,
          length: 0,
          weight: 0,
          aquatic: false,
          flying: false,
        });
      })
      .catch(error => console.error('Error creating dinosaur:', error));
  };

  // Обновление (UPDATE)
  const handleUpdate = () => {
    axios.put(`${API_BASE_URL}/${editingDino.id}`, editingDino)
      .then(response => {
        setDinosaurus(dinosaurus.map(dino => dino.id === editingDino.id ? response.data : dino));
        setEditingDino(null);
      })
      .catch(error => console.error('Error updating dinosaur:', error));
  };

  // Удаление (DELETE)
  const handleDelete = (id) => {
    axios.delete(`${API_BASE_URL}/${id}`)
      .then(() => setDinosaurus(dinosaurus.filter(dino => dino.id !== id)))
      .catch(error => console.error('Error deleting dinosaur:', error));
  };

  return (
    <div className="container mt-4">
      <h1 className="mb-4">Dinosaurus CRUD App</h1>

      {/* Форма для создания */}
      <h2>Create Dinosaur</h2>
      <div className="mb-3">
        <input
          type="text"
          className="form-control mb-2"
          placeholder="Species"
          value={newDino.species}
          onChange={e => setNewDino({ ...newDino, species: e.target.value })}
        />
        <input
          type="text"
          className="form-control mb-2"
          placeholder="Types"
          value={newDino.types}
          onChange={e => setNewDino({ ...newDino, types: e.target.value })}
        />
        <input
          type="number"
          className="form-control mb-2"
          placeholder="Height (m)"
          value={newDino.height}
          onChange={e => setNewDino({ ...newDino, height: parseFloat(e.target.value) })}
        />
        <input
          type="number"
          className="form-control mb-2"
          placeholder="Length (m)"
          value={newDino.length}
          onChange={e => setNewDino({ ...newDino, length: parseFloat(e.target.value) })}
        />
        <input
          type="number"
          className="form-control mb-2"
          placeholder="Weight (kg)"
          value={newDino.weight}
          onChange={e => setNewDino({ ...newDino, weight: parseFloat(e.target.value) })}
        />
        <div className="form-check mb-2">
          <input
            type="checkbox"
            className="form-check-input"
            checked={newDino.aquatic}
            onChange={e => setNewDino({ ...newDino, aquatic: e.target.checked })}
          />
          <label className="form-check-label">Aquatic</label>
        </div>
        <div className="form-check mb-2">
          <input
            type="checkbox"
            className="form-check-input"
            checked={newDino.flying}
            onChange={e => setNewDino({ ...newDino, flying: e.target.checked })}
          />
          <label className="form-check-label">Flying</label>
        </div>
        <button className="btn btn-primary" onClick={handleCreate}>Create</button>
      </div>

      {/* Список dinosaurs */}
      <h2>Dinosaur List</h2>
      <table className="table table-striped">
        <thead>
          <tr>
            <th>ID</th>
            <th>Species</th>
            <th>Types</th>
            <th>Height (m)</th>
            <th>Length (m)</th>
            <th>Weight (kg)</th>
            <th>Aquatic</th>
            <th>Flying</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {dinosaurus.map(dino => (
            <tr key={dino.id}>
              <td>{dino.id}</td>
              <td>{dino.species}</td>
              <td>{dino.types}</td>
              <td>{dino.height}</td>
              <td>{dino.length}</td>
              <td>{dino.weight}</td>
              <td>{dino.aquatic ? 'Yes' : 'No'}</td>
              <td>{dino.flying ? 'Yes' : 'No'}</td>
              <td>
                <button
                  className="btn btn-warning btn-sm me-2"
                  onClick={() => setEditingDino(dino)}
                >
                  Edit
                </button>
                <button
                  className="btn btn-danger btn-sm"
                  onClick={() => handleDelete(dino.id)}
                >
                  Delete
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>

      {/* Форма для редактирования */}
      {editingDino && (
        <div className="mt-4">
          <h2>Edit Dinosaur</h2>
          <div className="mb-3">
            <input
              type="text"
              className="form-control mb-2"
              value={editingDino.species}
              onChange={e => setEditingDino({ ...editingDino, species: e.target.value })}
            />
            <input
              type="text"
              className="form-control mb-2"
              value={editingDino.types}
              onChange={e => setEditingDino({ ...editingDino, types: e.target.value })}
            />
            <input
              type="number"
              className="form-control mb-2"
              value={editingDino.height}
              onChange={e => setEditingDino({ ...editingDino, height: parseFloat(e.target.value) })}
            />
            <input
              type="number"
              className="form-control mb-2"
              value={editingDino.length}
              onChange={e => setEditingDino({ ...editingDino, length: parseFloat(e.target.value) })}
            />
            <input
              type="number"
              className="form-control mb-2"
              value={editingDino.weight}
              onChange={e => setEditingDino({ ...editingDino, weight: parseFloat(e.target.value) })}
            />
            <div className="form-check mb-2">
              <input
                type="checkbox"
                className="form-check-input"
                checked={editingDino.aquatic}
                onChange={e => setEditingDino({ ...editingDino, aquatic: e.target.checked })}
              />
              <label className="form-check-label">Aquatic</label>
            </div>
            <div className="form-check mb-2">
              <input
                type="checkbox"
                className="form-check-input"
                checked={editingDino.flying}
                onChange={e => setEditingDino({ ...editingDino, flying: e.target.checked })}
              />
              <label className="form-check-label">Flying</label>
            </div>
            <button className="btn btn-primary" onClick={handleUpdate}>Update</button>
          </div>
        </div>
      )}
    </div>
  );
}

export default App;
