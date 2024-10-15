function showSuggestions() {
    const query = document.getElementById('search').value;
    if (query.length < 2) {
        document.getElementById('suggestions').innerHTML = '';
        return;
    }
    fetch(`http://localhost:8080/search?title=${encodeURIComponent(query)}`)
        .then(response => response.json())
        .then(data => {
            const suggestionsDiv = document.getElementById('suggestions');
            suggestionsDiv.innerHTML = '';
            data.forEach(film => {
                const suggestion = document.createElement('div');
                suggestion.textContent = film.original_title;
                suggestion.onclick = () => {
                    document.getElementById('search').value = film.original_title;
                    document.getElementById('suggestions').innerHTML = '';
                    displayFilmDetails(film);
                };
                suggestionsDiv.appendChild(suggestion);
            });
        })
        .catch(error => {
            console.error('エラーが発生しました:', error);
        });
}

function displayFilmDetails(film) {
    const filmsDiv = document.getElementById('films');
    filmsDiv.innerHTML = `
        <div class="film">
            <h2>${film.title} (${film.original_title})</h2>
            <p><strong>監督:</strong> ${film.director}</p>
            <p><strong>プロデューサー:</strong> ${film.producer}</p>
            <p><strong>公開年:</strong> ${film.release_date}</p>
            <p>${film.description}</p>
            <button onclick="fetchCharacters('${film.original_title}')">キャラクターを表示</button>
            <div id="characters"></div>
        </div>
    `;
}

function fetchCharacters(title) {
    fetch(`http://localhost:8080/search/characters?title=${encodeURIComponent(title)}`)
        .then(response => response.json())
        .then(data => {
            const charactersDiv = document.getElementById('characters');
            charactersDiv.innerHTML = '<h3>キャラクター一覧:</h3>';
            data.forEach(character => {
                const characterItem = document.createElement('p');
                characterItem.textContent = character.name;
                charactersDiv.appendChild(characterItem);
            });
        })
        .catch(error => {
            console.error('キャラクター情報の取得に失敗しました:', error);
        });
}
