import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { Container, Typography, Button, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Paper } from '@mui/material';
import { Link } from 'react-router-dom';

const DinosaurList = () => {
  const [dinosaurs, setDinosaurs] = useState([]);

  useEffect(() => {
    const fetchData = async () => {
      const token = localStorage.getItem('token');
      if (!token) {
        window.location.href = '/login';
        return;
      }
      try {
        const res = await axios.get('/api/dinosaurus', {
          headers: { Authorization: token }
        });
        setDinosaurs(res.data);
      } catch (err) {
        console.error('Error fetching dinosaurs:', err.response?.data || err.message);
        alert(`Error fetching data: ${err.response?.data?.error || err.message}`);
      }
    };
    fetchData();
  }, []);

  const handleDelete = async (id) => {
    const token = localStorage.getItem('token');
    try {
      await axios.delete(`/api/dinosaurus/${id}`, {
        headers: { Authorization: token }
      });
      setDinosaurs(dinosaurs.filter(d => d.id !== id));
    } catch (err) {
      alert('Error deleting dinosaur');
    }
  };

  return (
    <Container maxWidth="lg" sx={{ mt: 4 }}>
      <Typography variant="h4" gutterBottom>
        Dinosaurs
      </Typography>
      <Button variant="contained" color="primary" component={Link} to="/add" sx={{ mb: 2 }}>
        Add New Dinosaur
      </Button>
      <TableContainer component={Paper}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Species</TableCell>
              <TableCell>Types</TableCell>
              <TableCell>Height (m)</TableCell>
              <TableCell>Length (m)</TableCell>
              <TableCell>Weight (kg)</TableCell>
              <TableCell>Aquatic</TableCell>
              <TableCell>Flying</TableCell>
              <TableCell>Actions</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {dinosaurs.map(dino => (
              <TableRow key={dino.id}>
                <TableCell>{dino.species}</TableCell>
                <TableCell>{dino.types}</TableCell>
                <TableCell>{dino.height}</TableCell>
                <TableCell>{dino.length}</TableCell>
                <TableCell>{dino.weight}</TableCell>
                <TableCell>{dino.aquatic ? 'Yes' : 'No'}</TableCell>
                <TableCell>{dino.flying ? 'Yes' : 'No'}</TableCell>
                <TableCell>
                  <Button component={Link} to={`/edit/${dino.id}`} color="primary" sx={{ mr: 1 }}>
                    Edit
                  </Button>
                  <Button color="error" onClick={() => handleDelete(dino.id)}>
                    Delete
                  </Button>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    </Container>
  );
};

export default DinosaurList;