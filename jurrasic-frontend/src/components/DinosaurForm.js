import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import axios from 'axios';
import { Container, Typography, TextField, Checkbox, FormControlLabel, Button, Box } from '@mui/material';

const DinosaurForm = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const [dino, setDino] = useState({
    species: '', types: '', height: '', length: '', weight: '', aquatic: false, flying: false
  });

  useEffect(() => {
    if (id) {
      const fetchDino = async () => {
        const token = localStorage.getItem('token');
        try {
          const res = await axios.get(`/api/dinosaurus/${id}`, {
            headers: { Authorization: token }
          });
          setDino(res.data);
        } catch (err) {
          alert('Error fetching dinosaur');
        }
      };
      fetchDino();
    }
  }, [id]);

  const handleChange = (e) => {
    const { name, value, type, checked } = e.target;
    setDino({ ...dino, [name]: type === 'checkbox' ? checked : value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    const token = localStorage.getItem('token');
    // Преобразуем строковые значения в числа
    const payload = {
      ...dino,
      height: parseFloat(dino.height) || 0,
      length: parseFloat(dino.length) || 0,
      weight: parseFloat(dino.weight) || 0
    };
    try {
      if (id) {
        await axios.put(`/api/dinosaurus/${id}`, payload, {
          headers: { Authorization: token }
        });
      } else {
        await axios.post('/api/dinosaurus', payload, {
          headers: { Authorization: token }
        });
      }
      navigate('/');
    } catch (err) {
      console.error('Error saving dinosaur:', err.response?.data || err.message);
      alert(`Error saving dinosaur: ${err.response?.data?.error || err.message}`);
    }
  };

  return (
    <Container maxWidth="sm" sx={{ mt: 4 }}>
      <Typography variant="h4" gutterBottom>
        {id ? 'Edit Dinosaur' : 'Add Dinosaur'}
      </Typography>
      <Box component="form" onSubmit={handleSubmit} sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
        <TextField
          label="Species"
          name="species"
          value={dino.species}
          onChange={handleChange}
          fullWidth
          required
        />
        <TextField
          label="Types"
          name="types"
          value={dino.types}
          onChange={handleChange}
          fullWidth
          required
        />
        <TextField
          label="Height (m)"
          name="height"
          type="number"
          value={dino.height}
          onChange={handleChange}
          fullWidth
          required
          inputProps={{ step: '0.1' }}
        />
        <TextField
          label="Length (m)"
          name="length"
          type="number"
          value={dino.length}
          onChange={handleChange}
          fullWidth
          required
          inputProps={{ step: '0.1' }}
        />
        <TextField
          label="Weight (kg)"
          name="weight"
          type="number"
          value={dino.weight}
          onChange={handleChange}
          fullWidth
          required
          inputProps={{ step: '0.1' }}
        />
        <FormControlLabel
          control={<Checkbox name="aquatic" checked={dino.aquatic} onChange={handleChange} />}
          label="Aquatic"
        />
        <FormControlLabel
          control={<Checkbox name="flying" checked={dino.flying} onChange={handleChange} />}
          label="Flying"
        />
        <Button type="submit" variant="contained" color="primary">
          Save
        </Button>
      </Box>
    </Container>
  );
};

export default DinosaurForm;