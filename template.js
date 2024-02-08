import React, { useState } from 'react';
import './App.css';

function App() {
    const [activeSection, setActiveSection] = useState('inputForm');

    const changeSection = (sectionId) => {
        setActiveSection(sectionId);
    }

    const handleSubmit = (event) => {
        event.preventDefault();

        const expressionValue = event.target.expression.value;

        const expressionData = {
            expression: expressionValue
        };

        const idempotencyToken = generateUniqueToken();

        fetch('/expression', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'X-Idempotency-Token': idempotencyToken
            },
            body: JSON.stringify(expressionData)
        })
            .then(response => {
                if (!response.ok) {
                    throw new Error(`HTTP error! Status: ${response.status}`);
                }
                return response.json();
            })
            .then(data => {
                console.log(data);
            })
            .catch(error => {
                console.error('There was a problem with the fetch operation: ', error);
            });
    }

    return (
        <div className="App">
            <nav>
                <ul>
                    <li><button onClick={() => changeSection('inputForm')}>Input Expression</button></li>
                    <li><button onClick={() => changeSection('expressionsList')}>Expressions List</button></li>
                    <li><button onClick={() => changeSection('operationsList')}>Operations List</button></li>
                    <li><button onClick={() => changeSection('computationalResources')}>Computational Resources</button></li>
                </ul>
            </nav>

            <div id="content">
                <div id="inputForm" className={activeSection === 'inputForm' ? 'section active' : 'section'}>
                    <h2>Input Arithmetic Expression</h2>
                    <form id="expressionForm" onSubmit={handleSubmit}>
                        <label htmlFor="expression">Expression:</label>
                        <input type="text" id="expression" name="expression" />
                        <button type="submit">Submit</button>
                    </form>
                </div>

                <div id="expressionsList" className={activeSection === 'expressionsList' ? 'section active' : 'section'}>
                    <h2>Expressions List</h2>
                    <button id="getDataButtonExpressList">Get Data</button>
                    <ul id="expressionsListItems">
                        {/* Computational capabilities will be populated here */}
                    </ul>
                </div>
                {/* Остальные разделы также переведены в компоненты React */}
            </div>
        </div>
    );
}

// Функция для генерации уникального идентификатора
function generateUniqueToken() {
    return Math.random().toString(36).substring(2) + Date.now().toString(36);
}

export default App;